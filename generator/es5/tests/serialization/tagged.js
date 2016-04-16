$module('tagged', function () {
  var $static = this;
  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        somefield: SomeField,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'SomeField', {
        get: function () {
          if (this.$lazychecked['somefield']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['somefield'], $g.____testlib.basictypes.Integer, false, 'SomeField');
            this.$lazychecked['somefield'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['somefield'], $g.____testlib.basictypes.Integer);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['somefield'] = this.SomeField;
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
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['somefield'], right[BOXED_DATA_PROPERTY]['somefield'], $g.____testlib.basictypes.Integer));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['somefield'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['somefield'] = value;
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
            $g.tagged.SomeStruct.new($t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
            jsonString = $t.box('{"somefield":2}', $g.____testlib.basictypes.String);
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
