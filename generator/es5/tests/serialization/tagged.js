$module('tagged', function () {
  var $static = this;
  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.SomeField = SomeField;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['somefield'] = this.SomeField;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['somefield'], right.$data['somefield'], $g.____testlib.basictypes.Integer));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['somefield'], $g.____testlib.basictypes.Integer, false, 'SomeField');
        }
        if (this.$data['somefield'] != null) {
          return $t.box(this.$data['somefield'], $g.____testlib.basictypes.Integer);
        }
        return this.$data['somefield'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['somefield'] = $t.unbox(val);
          return;
        }
        this.$data['somefield'] = val;
      },
    });
  });

  $static.TEST = function () {
    var jsonString;
    var s;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.tagged.SomeStruct.new($t.nominalwrap(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
            jsonString = $t.nominalwrap('{"somefield":2}', $g.____testlib.basictypes.String);
            s.Stringify($g.____testlib.basictypes.JSON)().then(function ($result0) {
              return $g.____testlib.basictypes.String.$equals($result0, jsonString).then(function ($result1) {
                $result = $result1;
                $state.current = 2;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
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
