$module('nominal', function () {
  var $static = this;
  this.$type('CoolBool', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $global.Boolean;
    };
    this.$typesig = function () {
      return $t.createtypesig();
    };
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (someField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        someField: someField,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'someField', {
        get: function () {
          if (this.$lazychecked['someField']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['someField'], $g.____testlib.basictypes.Boolean, false, 'someField');
            this.$lazychecked['someField'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['someField'], $g.nominal.CoolBool);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['someField'] = this.someField;
        return $promise.resolve($t.box(mapped, $g.____testlib.basictypes.Mapping($t.any)));
      };
      return instance;
    };
    $instance.Mapping = function () {
      return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['someField'], right[BOXED_DATA_PROPERTY]['someField'], $g.nominal.CoolBool));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'someField', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['someField'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['someField'] = value;
      },
    });
    this.$typesig = function () {
      return $t.createtypesig(['someField', 5, $g.nominal.CoolBool.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.nominal.SomeStruct).$typeref()], ['Parse', 1, $g.____testlib.basictypes.Function($g.nominal.SomeStruct).$typeref()], ['equals', 4, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['Stringify', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()], ['Mapping', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Mapping($t.any)).$typeref()], ['Clone', 2, $g.____testlib.basictypes.Function($g.nominal.SomeStruct).$typeref()], ['String', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.String).$typeref()]);
    };
  });

  $static.TEST = function () {
    var c;
    var s;
    var s2;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            c = $t.box($t.box(true, $g.____testlib.basictypes.Boolean), $g.nominal.CoolBool);
            $g.nominal.SomeStruct.new(c).then(function ($result0) {
              $temp0 = $result0;
              $result = ($temp0, $temp0);
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            s = $result;
            $g.nominal.SomeStruct.Parse($g.____testlib.basictypes.JSON)($t.box('{"someField": true}', $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            s2 = $result;
            $promise.resolve(s2.someField).then(function ($result0) {
              $result = $t.box($result0 && s.someField, $g.____testlib.basictypes.Boolean);
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
