$module('nominal', function () {
  var $static = this;
  this.$type('CoolBool', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.new = function ($wrapped) {
      var instance = new this();
      instance.$wrapped = $wrapped;
      return instance;
    };
    this.$box = function (data) {
      var instance = new this();
      instance.$wrapped = data;
      return instance;
    };
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (someField) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.someField = someField;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['someField'] = this.someField;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['someField'], right.$data['someField'], $g.nominal.CoolBool));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'someField', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['someField'], $g.____testlib.basictypes.Boolean, false, 'someField');
        }
        if (this.$data['someField'] != null) {
          return $t.box(this.$data['someField'], $g.nominal.CoolBool);
        }
        return this.$data['someField'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['someField'] = $t.unbox(val);
          return;
        }
        this.$data['someField'] = val;
      },
    });
  });

  $static.TEST = function () {
    var c;
    var s;
    var s2;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            c = $t.nominalwrap($t.nominalwrap(true, $g.____testlib.basictypes.Boolean), $g.nominal.CoolBool);
            $g.nominal.SomeStruct.new(c).then(function ($result0) {
              $temp0 = $result0;
              $result = ($temp0, $temp0);
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            s = $result;
            $g.nominal.SomeStruct.Parse($g.____testlib.basictypes.JSON)($t.nominalwrap('{"someField": true}', $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            s2 = $result;
            $promise.resolve(s2.someField).then(function ($result0) {
              $result = $t.nominalwrap($result0 && s.someField, $g.____testlib.basictypes.Boolean);
              $state.current = 3;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $state.resolve($result);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
